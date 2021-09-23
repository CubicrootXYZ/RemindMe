"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.defaultPlugin = exports.builtInConfigs = void 0;
const recommended_1 = require("./recommended");
const all_1 = require("./all");
const minimal_1 = require("./minimal");
const builtinRules = require("../rules/builtin");
exports.builtInConfigs = {
    recommended: recommended_1.default,
    minimal: minimal_1.default,
    all: all_1.default,
    'redocly-registry': {
        decorators: { 'registry-dependencies': 'on' }
    }
};
exports.defaultPlugin = {
    id: '',
    rules: builtinRules.rules,
    preprocessors: builtinRules.preprocessors,
    decorators: builtinRules.decorators,
    configs: exports.builtInConfigs,
};
