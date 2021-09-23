"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ValidContentExamples = void 0;
const ajv_1 = require("../ajv");
const ref_utils_1 = require("../../ref-utils");
const ValidContentExamples = (opts) => {
    var _a;
    const disallowAdditionalProperties = (_a = opts.disallowAdditionalProperties) !== null && _a !== void 0 ? _a : true;
    return {
        MediaType: {
            leave(mediaType, { report, location, resolve }) {
                if (!mediaType.schema)
                    return;
                if (mediaType.example) {
                    validateExample(mediaType.example, location.child('example'));
                }
                else if (mediaType.examples) {
                    for (const exampleName of Object.keys(mediaType.examples)) {
                        let example = mediaType.examples[exampleName];
                        let dataLoc = location.child(['examples', exampleName, 'value']);
                        if (ref_utils_1.isRef(example)) {
                            const resolved = resolve(example);
                            if (!resolved.location)
                                continue;
                            dataLoc = resolved.location.child('value');
                            example = resolved.node;
                        }
                        validateExample(example.value, dataLoc);
                    }
                }
                function validateExample(example, dataLoc) {
                    try {
                        const { valid, errors } = ajv_1.validateJsonSchema(example, mediaType.schema, location.child('schema'), dataLoc.pointer, resolve, disallowAdditionalProperties);
                        if (!valid) {
                            for (let error of errors) {
                                report({
                                    message: `Example value must conform to the schema: ${error.message}.`,
                                    location: Object.assign(Object.assign({}, new ref_utils_1.Location(dataLoc.source, error.instancePath)), { reportOnKey: error.keyword === 'additionalProperties' }),
                                    from: location,
                                    suggest: error.suggest,
                                });
                            }
                        }
                    }
                    catch (e) {
                        report({
                            message: `Example validation errored: ${e.message}.`,
                            location: location.child('schema'),
                            from: location
                        });
                    }
                }
            },
        },
    };
};
exports.ValidContentExamples = ValidContentExamples;
