"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.validateJsonSchema = exports.releaseAjvInstance = void 0;
const ajv_1 = require("@redocly/ajv");
const ref_utils_1 = require("../ref-utils");
let ajvInstance = null;
function releaseAjvInstance() {
    ajvInstance = null;
}
exports.releaseAjvInstance = releaseAjvInstance;
function getAjv(resolve, disallowAdditionalProperties) {
    if (!ajvInstance) {
        ajvInstance = new ajv_1.default({
            schemaId: '$id',
            meta: true,
            allErrors: true,
            strictSchema: false,
            inlineRefs: false,
            validateSchema: false,
            discriminator: true,
            allowUnionTypes: true,
            validateFormats: false,
            defaultAdditionalProperties: !disallowAdditionalProperties,
            loadSchemaSync(base, $ref) {
                const resolvedRef = resolve({ $ref }, base.split('#')[0]);
                if (!resolvedRef || !resolvedRef.location)
                    return undefined;
                return Object.assign({ $id: resolvedRef.location.absolutePointer }, resolvedRef.node);
            },
            logger: false,
        });
    }
    return ajvInstance;
}
function getAjvValidator(schema, loc, resolve, disallowAdditionalProperties) {
    const ajv = getAjv(resolve, disallowAdditionalProperties);
    if (!ajv.getSchema(loc.absolutePointer)) {
        ajv.addSchema(Object.assign({ $id: loc.absolutePointer }, schema), loc.absolutePointer);
    }
    return ajv.getSchema(loc.absolutePointer);
}
function validateJsonSchema(data, schema, schemaLoc, instancePath, resolve, disallowAdditionalProperties) {
    const validate = getAjvValidator(schema, schemaLoc, resolve, disallowAdditionalProperties);
    if (!validate)
        return { valid: true, errors: [] }; // unresolved refs are reported
    const valid = validate(data, {
        instancePath,
        parentData: { fake: {} },
        parentDataProperty: 'fake',
        rootData: {},
        dynamicAnchors: {},
    });
    return {
        valid: !!valid,
        errors: (validate.errors || []).map(beatifyErrorMessage),
    };
    function beatifyErrorMessage(error) {
        let message = error.message;
        let suggest = error.keyword === 'enum' ? error.params.allowedValues : undefined;
        if (suggest) {
            message += ` ${suggest.map((e) => `"${e}"`).join(', ')}`;
        }
        if (error.keyword === 'type') {
            message = `type ${message}`;
        }
        const relativePath = error.instancePath.substring(instancePath.length + 1);
        const propName = relativePath.substring(relativePath.lastIndexOf('/') + 1);
        if (propName) {
            message = `\`${propName}\` property ${message}`;
        }
        if (error.keyword === 'additionalProperties') {
            const property = error.params.additionalProperty;
            message = `${message} \`${property}\``;
            error.instancePath += '/' + ref_utils_1.escapePointer(property);
        }
        return Object.assign(Object.assign({}, error), { message,
            suggest });
    }
}
exports.validateJsonSchema = validateJsonSchema;
