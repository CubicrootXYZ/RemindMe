"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getSuggest = exports.validateDefinedAndNonEmpty = exports.fieldNonEmpty = exports.missingRequiredField = exports.matchesJsonSchemaType = exports.oasTypeOf = void 0;
const levenshtein = require("js-levenshtein");
function oasTypeOf(value) {
    if (Array.isArray(value)) {
        return 'array';
    }
    else if (value === null) {
        return 'null';
    }
    else {
        return typeof value;
    }
}
exports.oasTypeOf = oasTypeOf;
/**
 * Checks if value matches specified JSON schema type
 *
 * @param {*} value - value to check
 * @param {JSONSchemaType} type - JSON Schema type
 * @returns boolean
 */
function matchesJsonSchemaType(value, type, nullable) {
    if (nullable && value === null) {
        return value === null;
    }
    switch (type) {
        case 'array':
            return Array.isArray(value);
        case 'object':
            return typeof value === 'object' && value !== null && !Array.isArray(value);
        case 'null':
            return value === null;
        case 'integer':
            return Number.isInteger(value);
        default:
            return typeof value === type;
    }
}
exports.matchesJsonSchemaType = matchesJsonSchemaType;
function missingRequiredField(type, field) {
    return `${type} object should contain \`${field}\` field.`;
}
exports.missingRequiredField = missingRequiredField;
function fieldNonEmpty(type, field) {
    return `${type} object \`${field}\` must be non-empty string.`;
}
exports.fieldNonEmpty = fieldNonEmpty;
function validateDefinedAndNonEmpty(fieldName, value, ctx) {
    if (typeof value !== 'object') {
        return;
    }
    if (value[fieldName] === undefined) {
        ctx.report({
            message: missingRequiredField(ctx.type.name, fieldName),
            location: ctx.location.child([fieldName]).key(),
        });
    }
    else if (!value[fieldName]) {
        ctx.report({
            message: fieldNonEmpty(ctx.type.name, fieldName),
            location: ctx.location.child([fieldName]).key(),
        });
    }
}
exports.validateDefinedAndNonEmpty = validateDefinedAndNonEmpty;
function getSuggest(given, variants) {
    if (typeof given !== 'string' || !variants.length)
        return [];
    const distances = [];
    for (let i = 0; i < variants.length; i++) {
        const distance = levenshtein(given, variants[i]);
        if (distance < 4) {
            distances.push({ distance, variant: variants[i] });
        }
    }
    distances.sort((a, b) => a.distance - b.distance);
    // if (bestMatch.distance <= 4) return bestMatch.string;
    return distances.map((d) => d.variant);
}
exports.getSuggest = getSuggest;
