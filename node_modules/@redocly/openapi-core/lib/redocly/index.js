"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.RedoclyClient = void 0;
const fs_1 = require("fs");
const path_1 = require("path");
const os_1 = require("os");
const colorette_1 = require("colorette");
const query_1 = require("./query");
const TOKEN_FILENAME = '.redocly-config.json';
class RedoclyClient {
    constructor() {
        this.loadToken();
    }
    hasToken() {
        return !!this.accessToken;
    }
    loadToken() {
        if (process.env.REDOCLY_AUTHORIZATION) {
            this.accessToken = process.env.REDOCLY_AUTHORIZATION;
            return;
        }
        const credentialsPath = path_1.resolve(os_1.homedir(), TOKEN_FILENAME);
        if (fs_1.existsSync(credentialsPath)) {
            const credentials = JSON.parse(fs_1.readFileSync(credentialsPath, 'utf-8'));
            this.accessToken = credentials && credentials.token;
        }
    }
    isAuthorizedWithRedocly() {
        return __awaiter(this, void 0, void 0, function* () {
            return this.hasToken() && !!(yield this.getAuthorizationHeader());
        });
    }
    verifyToken(accessToken, verbose = false) {
        return __awaiter(this, void 0, void 0, function* () {
            if (!accessToken)
                return false;
            const authDetails = yield RedoclyClient.authorize(accessToken, { verbose });
            if (!authDetails)
                return false;
            return true;
        });
    }
    getAuthorizationHeader() {
        return __awaiter(this, void 0, void 0, function* () {
            // print this only if there is token but invalid
            if (this.accessToken && !(yield this.verifyToken(this.accessToken))) {
                process.stderr.write(`${colorette_1.yellow('Warning:')} invalid Redocly API key. Use "npx @redocly/openapi-cli login" to provide your API key\n`);
                return undefined;
            }
            return this.accessToken;
        });
    }
    login(accessToken, verbose = false) {
        return __awaiter(this, void 0, void 0, function* () {
            const credentialsPath = path_1.resolve(os_1.homedir(), TOKEN_FILENAME);
            process.stdout.write(colorette_1.gray('\n  Logging in...\n'));
            const authorized = yield this.verifyToken(accessToken, verbose);
            if (!authorized) {
                process.stdout.write(colorette_1.red('Authorization failed. Please check if you entered a valid API key.\n'));
                process.exit(1);
            }
            this.accessToken = accessToken;
            const credentials = {
                token: accessToken,
            };
            fs_1.writeFileSync(credentialsPath, JSON.stringify(credentials, null, 2));
            process.stdout.write(colorette_1.green('  Authorization confirmed. ✅\n\n'));
        });
    }
    logout() {
        const credentialsPath = path_1.resolve(os_1.homedir(), TOKEN_FILENAME);
        if (fs_1.existsSync(credentialsPath)) {
            fs_1.unlinkSync(credentialsPath);
        }
        process.stdout.write('Logged out from the Redocly account. ✋\n');
    }
    query(queryString, parameters = {}, headers = {}) {
        return __awaiter(this, void 0, void 0, function* () {
            return query_1.query(queryString, parameters, Object.assign({ Authorization: this.accessToken }, headers));
        });
    }
    static authorize(accessToken, options) {
        return __awaiter(this, void 0, void 0, function* () {
            const { queryName = '', verbose = false } = options;
            try {
                const queryStr = `query ${queryName}{ viewer { id } }`;
                return yield query_1.query(queryStr, {}, { Authorization: accessToken });
            }
            catch (e) {
                if (verbose)
                    console.log(e);
                return null;
            }
        });
    }
    updateDependencies(dependencies) {
        return __awaiter(this, void 0, void 0, function* () {
            const definitionId = process.env.DEFINITION;
            const versionId = process.env.DEFINITION;
            const branchId = process.env.BRANCH;
            if (!definitionId || !versionId || !branchId)
                return;
            yield this.query(`
    mutation UpdateBranchDependenciesFromURLs(
      $urls: [String!]!
      $definitionId: Int!
      $versionId: Int!
      $branchId: Int!
    ) {
      updateBranchDependenciesFromURLs(
        definitionId: $definitionId
        versionId: $versionId
        branchId: $branchId
        urls: $urls
      ) {
        branchName
      }
    }
    `, {
                urls: dependencies || [],
                definitionId: parseInt(definitionId, 10),
                versionId: parseInt(versionId, 10),
                branchId: parseInt(branchId, 10),
            });
        });
    }
    updateDefinitionVersion(definitionId, versionId, updatePatch) {
        return this.query(`
      mutation UpdateDefinitionVersion($definitionId: Int!, $versionId: Int!, $updatePatch: DefinitionVersionPatch!) {
        updateDefinitionVersionByDefinitionIdAndId(input: {definitionId: $definitionId, id: $versionId, patch: $updatePatch}) {
          definitionVersion {
            ...VersionDetails
            __typename
          }
          __typename
        }
      }

      fragment VersionDetails on DefinitionVersion {
        id
        nodeId
        uuid
        definitionId
        name
        description
        sourceType
        source
        registryAccess
        __typename
      }
    `, {
            definitionId,
            versionId,
            updatePatch,
        });
    }
    getOrganizationId(organizationId) {
        return this.query(`
      query ($organizationId: String!) {
        organizationById(id: $organizationId) {
          id
        }
      }
    `, {
            organizationId
        });
    }
    getDefinitionByName(name, organizationId) {
        return this.query(`
      query ($name: String!, $organizationId: String!) {
        definition: definitionByOrganizationIdAndName(name: $name, organizationId: $organizationId) {
          id
        }
      }
    `, {
            name,
            organizationId
        });
    }
    createDefinition(organizationId, name) {
        return this.query(`
      mutation CreateDefinition($organizationId: String!, $name: String!) {
        def: createDefinition(input: {organizationId: $organizationId, name: $name }) {
          definition {
            id
            nodeId
            name
          }
        }
      }
    `, {
            organizationId,
            name
        });
    }
    createDefinitionVersion(definitionId, name, sourceType, source) {
        return this.query(`
      mutation CreateVersion($definitionId: Int!, $name: String!, $sourceType: DvSourceType!, $source: JSON) {
        createDefinitionVersion(input: {definitionId: $definitionId, name: $name, sourceType: $sourceType, source: $source }) {
          definitionVersion {
            id
          }
        }
      }
    `, {
            definitionId,
            name,
            sourceType,
            source
        });
    }
    getSignedUrl(organizationId, filesHash, fileName) {
        return this.query(`
      query ($organizationId: String!, $filesHash: String!, $fileName: String!) {
        signFileUploadCLI(organizationId: $organizationId, filesHash: $filesHash, fileName: $fileName) {
          signedFileUrl
          uploadedFilePath
        }
      }
    `, {
            organizationId,
            filesHash,
            fileName
        });
    }
    getDefinitionVersion(organizationId, definitionName, versionName) {
        return this.query(`
      query ($organizationId: String!, $definitionName: String!, $versionName: String!) {
        version: definitionVersionByOrganizationDefinitionAndName(organizationId: $organizationId, definitionName: $definitionName, versionName: $versionName) {
          id
          definitionId
          defaultBranch {
            name
          }
        }
      }
    `, {
            organizationId,
            definitionName,
            versionName
        });
    }
    static isRegistryURL(link) {
        const domain = process.env.REDOCLY_DOMAIN || 'redoc.ly';
        if (!link.startsWith(`https://api.${domain}/registry/`))
            return false;
        const registryPath = link.replace(`https://api.${domain}/registry/`, '');
        const pathParts = registryPath.split('/');
        // we can be sure, that there is job UUID present
        // (org, definition, version, bundle, branch, job, "openapi.yaml" 🤦‍♂️)
        // so skip this link.
        // FIXME
        if (pathParts.length === 7)
            return false;
        return true;
    }
}
exports.RedoclyClient = RedoclyClient;
