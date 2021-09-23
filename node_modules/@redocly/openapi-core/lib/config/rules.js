"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.initRules = void 0;
const utils_1 = require("../utils");
function initRules(rules, config, type, oasVersion) {
    return rules
        .flatMap((ruleset) => Object.keys(ruleset).map((ruleId) => {
        const rule = ruleset[ruleId];
        const ruleSettings = type === 'rules'
            ? config.getRuleSettings(ruleId, oasVersion)
            : type === 'preprocessors'
                ? config.getPreprocessorSettings(ruleId, oasVersion)
                : config.getDecoratorSettings(ruleId, oasVersion);
        if (ruleSettings.severity === 'off') {
            return undefined;
        }
        const visitor = rule(ruleSettings);
        return {
            severity: ruleSettings.severity,
            ruleId,
            visitor,
        };
    }))
        .filter(utils_1.notUndefined);
}
exports.initRules = initRules;
