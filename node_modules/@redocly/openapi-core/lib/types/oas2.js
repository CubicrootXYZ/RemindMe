"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Oas2Types = void 0;
const _1 = require(".");
const responseCodeRegexp = /^[0-9][0-9Xx]{2}$/;
const DefinitionRoot = {
    properties: {
        swagger: { type: 'string' },
        info: 'Info',
        host: { type: 'string' },
        basePath: { type: 'string' },
        schemes: { type: 'array', items: { type: 'string' } },
        consumes: { type: 'array', items: { type: 'string' } },
        produces: { type: 'array', items: { type: 'string' } },
        paths: 'PathMap',
        definitions: 'NamedSchemas',
        parameters: 'NamedParameters',
        responses: 'NamedResponses',
        securityDefinitions: 'NamedSecuritySchemes',
        security: _1.listOf('SecurityRequirement'),
        tags: _1.listOf('Tag'),
        externalDocs: 'ExternalDocs',
    },
    required: ['swagger', 'paths', 'info'],
};
const Info = {
    properties: {
        title: { type: 'string' },
        description: { type: 'string' },
        termsOfService: { type: 'string' },
        contact: 'Contact',
        license: 'License',
        version: { type: 'string' },
    },
    required: ['title', 'version'],
};
const Contact = {
    properties: {
        name: { type: 'string' },
        url: { type: 'string' },
        email: { type: 'string' },
    },
};
const License = {
    properties: {
        name: { type: 'string' },
        url: { type: 'string' },
    },
    required: ['name'],
};
const PathMap = {
    properties: {},
    additionalProperties: (_value, key) => key.startsWith('/') ? 'PathItem' : undefined,
};
const PathItem = {
    properties: {
        $ref: { type: 'string' },
        get: 'Operation',
        put: 'Operation',
        post: 'Operation',
        delete: 'Operation',
        options: 'Operation',
        head: 'Operation',
        patch: 'Operation',
        parameters: _1.listOf('Parameter'),
    },
};
const Operation = {
    properties: {
        tags: { type: 'array', items: { type: 'string' } },
        summary: {
            type: 'string',
        },
        description: { type: 'string' },
        externalDocs: 'ExternalDocs',
        operationId: { type: 'string' },
        consumes: { type: 'array', items: { type: 'string' } },
        produces: { type: 'array', items: { type: 'string' } },
        parameters: _1.listOf('Parameter'),
        responses: 'ResponsesMap',
        schemes: { type: 'array', items: { type: 'string' } },
        deprecated: { type: 'boolean' },
        security: _1.listOf('SecurityRequirement'),
        'x-codeSamples': _1.listOf('XCodeSample'),
        'x-code-samples': _1.listOf('XCodeSample'), // deprecated
    },
    required: ['responses'],
};
const XCodeSample = {
    properties: {
        lang: { type: 'string' },
        label: { type: 'string' },
        source: { type: 'string' },
    },
};
const ExternalDocs = {
    properties: {
        description: { type: 'string' },
        url: { type: 'string' },
    },
    required: ['url'],
};
const Parameter = {
    properties: {
        name: { type: 'string' },
        in: { type: 'string', enum: ['query', 'header', 'path', 'formData', 'body'] },
        description: { type: 'string' },
        required: { type: 'boolean' },
        schema: 'Schema',
        type: { type: 'string', enum: ['string', 'number', 'integer', 'boolean', 'array', 'file'] },
        format: { type: 'string' },
        allowEmptyValue: { type: 'boolean' },
        items: 'ParameterItems',
        collectionFormat: { type: 'string', enum: ['csv', 'ssv', 'tsv', 'pipes', 'multi'] },
        default: null,
        maximum: { type: 'integer' },
        exclusiveMaximum: { type: 'boolean' },
        minimum: { type: 'integer' },
        exclusiveMinimum: { type: 'boolean' },
        maxLength: { type: 'integer' },
        minLength: { type: 'integer' },
        pattern: { type: 'string' },
        maxItems: { type: 'integer' },
        minItems: { type: 'integer' },
        uniqueItems: { type: 'boolean' },
        enum: { type: 'array' },
        multipleOf: { type: 'number' },
    },
    required(value) {
        if (!value || !value.in) {
            return ['name', 'in'];
        }
        if (value.in === 'body') {
            return ['name', 'in', 'schema'];
        }
        else {
            if (value.type === 'array') {
                return ['name', 'in', 'type', 'items'];
            }
            else {
                return ['name', 'in', 'type'];
            }
        }
    },
};
const ParameterItems = {
    properties: {
        type: { type: 'string', enum: ['string', 'number', 'integer', 'boolean', 'array'] },
        format: { type: 'string' },
        items: 'ParameterItems',
        collectionFormat: { type: 'string', enum: ['csv', 'ssv', 'tsv', 'pipes', 'multi'] },
        default: null,
        maximum: { type: 'integer' },
        exclusiveMaximum: { type: 'boolean' },
        minimum: { type: 'integer' },
        exclusiveMinimum: { type: 'boolean' },
        maxLength: { type: 'integer' },
        minLength: { type: 'integer' },
        pattern: { type: 'string' },
        maxItems: { type: 'integer' },
        minItems: { type: 'integer' },
        uniqueItems: { type: 'boolean' },
        enum: { type: 'array' },
        multipleOf: { type: 'number' },
    },
    required(value) {
        if (value && value.type === 'array') {
            return ['type', 'items'];
        }
        else {
            return ['type'];
        }
    },
};
const ResponsesMap = {
    properties: {
        default: 'Response',
    },
    additionalProperties: (_v, key) => responseCodeRegexp.test(key) ? 'Response' : undefined,
};
const Response = {
    properties: {
        description: {
            type: 'string',
        },
        schema: 'Schema',
        headers: _1.mapOf('Header'),
        examples: 'Examples',
    },
    required: ['description'],
};
const Examples = {
    properties: {},
    additionalProperties: { isExample: true },
};
const Header = {
    properties: {
        description: { type: 'string' },
        type: { type: 'string', enum: ['string', 'number', 'integer', 'boolean', 'array'] },
        format: { type: 'string' },
        items: 'ParameterItems',
        collectionFormat: { type: 'string', enum: ['csv', 'ssv', 'tsv', 'pipes', 'multi'] },
        default: null,
        maximum: { type: 'integer' },
        exclusiveMaximum: { type: 'boolean' },
        minimum: { type: 'integer' },
        exclusiveMinimum: { type: 'boolean' },
        maxLength: { type: 'integer' },
        minLength: { type: 'integer' },
        pattern: { type: 'string' },
        maxItems: { type: 'integer' },
        minItems: { type: 'integer' },
        uniqueItems: { type: 'boolean' },
        enum: { type: 'array' },
        multipleOf: { type: 'number' },
    },
    required(value) {
        if (value && value.type === 'array') {
            return ['type', 'items'];
        }
        else {
            return ['type'];
        }
    },
};
const Tag = {
    properties: {
        name: { type: 'string' },
        description: { type: 'string' },
        externalDocs: 'ExternalDocs',
    },
    required: ['name'],
};
const Schema = {
    properties: {
        format: { type: 'string' },
        title: { type: 'string' },
        description: { type: 'string' },
        default: null,
        multipleOf: { type: 'number' },
        maximum: { type: 'number' },
        minimum: { type: 'number' },
        exclusiveMaximum: { type: 'boolean' },
        exclusiveMinimum: { type: 'boolean' },
        maxLength: { type: 'number' },
        minLength: { type: 'number' },
        pattern: { type: 'string' },
        maxItems: { type: 'number' },
        minItems: { type: 'number' },
        uniqueItems: { type: 'boolean' },
        maxProperties: { type: 'number' },
        minProperties: { type: 'number' },
        required: { type: 'array', items: { type: 'string' } },
        enum: { type: 'array' },
        type: {
            type: 'string',
            enum: ['object', 'array', 'string', 'number', 'integer', 'boolean', 'null'],
        },
        items: (value) => {
            if (Array.isArray(value)) {
                return _1.listOf('Schema');
            }
            else {
                return 'Schema';
            }
        },
        allOf: _1.listOf('Schema'),
        properties: 'SchemaProperties',
        additionalProperties: (value) => {
            if (typeof value === 'boolean') {
                return { type: 'boolean' };
            }
            else {
                return 'Schema';
            }
        },
        discriminator: { type: 'string' },
        readOnly: { type: 'boolean' },
        xml: 'Xml',
        externalDocs: 'ExternalDocs',
        example: { isExample: true },
    },
};
const SchemaProperties = {
    properties: {},
    additionalProperties: 'Schema',
};
const Xml = {
    properties: {
        name: { type: 'string' },
        namespace: { type: 'string' },
        prefix: { type: 'string' },
        attribute: { type: 'boolean' },
        wrapped: { type: 'boolean' },
    },
};
const SecurityScheme = {
    properties: {
        type: { enum: ['basic', 'apiKey', 'oauth2'] },
        description: { type: 'string' },
        name: { type: 'string' },
        in: { type: 'string', enum: ['query', 'header', 'cookie'] },
        flow: { enum: ['implicit', 'password', 'application', 'accessCode'] },
        authorizationUrl: { type: 'string' },
        tokenUrl: { type: 'string' },
        scopes: { type: 'object', additionalProperties: { type: 'string' } },
    },
    required(value) {
        if (!(value === null || value === void 0 ? void 0 : value.type)) {
            return ['type'];
        }
        if (value.type === 'apiKey') {
            return ['type', 'name', 'in'];
        }
        else if (value.type === 'http') {
            return ['type', 'scheme'];
        }
        else if (value.type === 'oauth2') {
            if (!(value === null || value === void 0 ? void 0 : value.flow)) {
                return ['type', 'flow'];
            }
            else if (value.flow === 'implicit') {
                return ['type', 'flow', 'authorizationUrl'];
            }
            else if (value.flow === 'accessCode') {
                return ['type', 'flow', 'authorizationUrl', 'tokenUrl'];
            }
            else if (value.flow === 'application') {
                return ['type', 'flow', 'tokenUrl'];
            }
            else if (value.flow === 'password') {
                return ['type', 'flow', 'tokenUrl'];
            }
            else {
                return ['type', 'flow'];
            }
        }
        return ['type'];
    },
};
const SecurityRequirement = {
    properties: {},
    additionalProperties: { type: 'array', items: { type: 'string' } },
};
exports.Oas2Types = {
    DefinitionRoot,
    Tag,
    ExternalDocs,
    SecurityRequirement,
    Info,
    Contact,
    License,
    PathMap,
    PathItem,
    Parameter,
    ParameterItems,
    Operation,
    Examples,
    Header,
    ResponsesMap,
    Response,
    Schema,
    Xml,
    SchemaProperties,
    NamedSchemas: _1.mapOf('Schema'),
    NamedResponses: _1.mapOf('Response'),
    NamedParameters: _1.mapOf('Parameter'),
    NamedSecuritySchemes: _1.mapOf('SecurityScheme'),
    SecurityScheme,
    XCodeSample,
};
