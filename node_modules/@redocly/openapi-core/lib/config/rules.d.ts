import { RuleSet, OasVersion } from '../oas-types';
import { LintConfig } from './config';
export declare function initRules<T extends Function, P extends RuleSet<T>>(rules: P[], config: LintConfig, type: 'rules' | 'preprocessors' | 'decorators', oasVersion: OasVersion): {
    severity: import("..").ProblemSeverity;
    ruleId: string;
    visitor: any;
}[];
