import { Document } from '../resolve';
import { Oas3RuleSet } from '../oas-types';
import { LintConfig, Plugin } from '../config/config';
export declare function parseYamlToDocument(body: string, absoluteRef?: string): Document;
export declare function makeConfigForRuleset(rules: Oas3RuleSet, plugin?: Partial<Plugin>): LintConfig;
