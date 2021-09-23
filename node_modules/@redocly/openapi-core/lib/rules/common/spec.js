"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.OasSpec = void 0;
const types_1 = require("../../types");
const utils_1 = require("../utils");
const ref_utils_1 = require("../../ref-utils");
const OasSpec = () => {
    return {
        any(node, { report, type, location, key, resolve, ignoreNextVisitorsOnNode }) {
            var _a, _b;
            const nodeType = utils_1.oasTypeOf(node);
            if (type.items) {
                if (nodeType !== 'array') {
                    report({
                        message: `Expected type \`${type.name}\` (array) but got \`${nodeType}\``,
                    });
                    ignoreNextVisitorsOnNode();
                }
                return;
            }
            else if (nodeType !== 'object') {
                report({
                    message: `Expected type \`${type.name}\` (object) but got \`${nodeType}\``,
                });
                ignoreNextVisitorsOnNode();
                return;
            }
            const required = typeof type.required === 'function' ? type.required(node, key) : type.required;
            for (let propName of required || []) {
                if (!node.hasOwnProperty(propName)) {
                    report({
                        message: `The field \`${propName}\` must be present on this level.`,
                        location: [{ reportOnKey: true }],
                    });
                }
            }
            const requiredOneOf = type.requiredOneOf || null;
            if (requiredOneOf) {
                let hasProperty = false;
                for (let propName of requiredOneOf || []) {
                    if (node.hasOwnProperty(propName)) {
                        hasProperty = true;
                    }
                }
                if (!hasProperty)
                    report({
                        message: 'Must contain at least one of the following fields: path, components, webhooks.',
                        location: [{ reportOnKey: true }],
                    });
            }
            for (const propName of Object.keys(node)) {
                const propLocation = location.child([propName]);
                let propValue = node[propName];
                let propType = type.properties[propName];
                if (propType === undefined)
                    propType = type.additionalProperties;
                if (typeof propType === 'function')
                    propType = propType(propValue, propName);
                if (types_1.isNamedType(propType)) {
                    continue; // do nothing for named schema, it is executed with the next any call
                }
                const propSchema = propType;
                const propValueType = utils_1.oasTypeOf(propValue);
                if (propSchema === undefined) {
                    if (propName.startsWith('x-'))
                        continue;
                    report({
                        message: `Property \`${propName}\` is not expected here.`,
                        suggest: utils_1.getSuggest(propName, Object.keys(type.properties)),
                        location: propLocation.key(),
                    });
                    continue;
                }
                if (propSchema === null) {
                    continue; // just defined, no validation
                }
                if (propSchema.resolvable !== false && ref_utils_1.isRef(propValue)) {
                    propValue = resolve(propValue).node;
                }
                if (propSchema.enum) {
                    if (!propSchema.enum.includes(propValue)) {
                        report({
                            location: propLocation,
                            message: `\`${propName}\` can be one of the following only: ${propSchema.enum
                                .map((i) => `"${i}"`)
                                .join(', ')}.`,
                            suggest: utils_1.getSuggest(propValue, propSchema.enum),
                        });
                    }
                }
                else if (propSchema.type && !utils_1.matchesJsonSchemaType(propValue, propSchema.type, false)) {
                    report({
                        message: `Expected type \`${propSchema.type}\` but got \`${propValueType}\`.`,
                        location: propLocation,
                    });
                }
                else if (propValueType === 'array' && ((_a = propSchema.items) === null || _a === void 0 ? void 0 : _a.type)) {
                    const itemsType = (_b = propSchema.items) === null || _b === void 0 ? void 0 : _b.type;
                    for (let i = 0; i < propValue.length; i++) {
                        const item = propValue[i];
                        if (!utils_1.matchesJsonSchemaType(item, itemsType, false)) {
                            report({
                                message: `Expected type \`${itemsType}\` but got \`${utils_1.oasTypeOf(item)}\`.`,
                                location: propLocation.child([i]),
                            });
                        }
                    }
                }
                if (typeof propSchema.minimum === 'number') {
                    if (propSchema.minimum > node[propName]) {
                        report({
                            message: `The value of the ${propName} field must be greater than or equal to ${propSchema.minimum}`,
                            location: location.child([propName]),
                        });
                    }
                }
            }
        },
    };
};
exports.OasSpec = OasSpec;
