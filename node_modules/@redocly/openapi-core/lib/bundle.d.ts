import { BaseResolver, Document } from './resolve';
import { Oas3Rule } from './visitors';
import { NormalizedNodeType, NodeType } from './types';
import type { Config, LintConfig } from './config/config';
export declare type Oas3RuleSet = Record<string, Oas3Rule>;
export declare enum OasVersion {
    Version2 = "oas2",
    Version3_0 = "oas3_0",
    Version3_1 = "oas3_1"
}
export declare function bundle(opts: {
    ref?: string;
    doc?: Document;
    externalRefResolver?: BaseResolver;
    config: Config;
    dereference?: boolean;
    base?: string;
}): Promise<{
    bundle: Document;
    problems: import("./walk").NormalizedProblem[];
    fileDependencies: Set<string>;
    rootType: NormalizedNodeType;
    refTypes: Map<string, NormalizedNodeType> | undefined;
}>;
export declare function bundleDocument(opts: {
    document: Document;
    config: LintConfig;
    customTypes?: Record<string, NodeType>;
    externalRefResolver: BaseResolver;
    dereference?: boolean;
}): Promise<{
    bundle: Document;
    problems: import("./walk").NormalizedProblem[];
    fileDependencies: Set<string>;
    rootType: NormalizedNodeType;
    refTypes: Map<string, NormalizedNodeType> | undefined;
}>;
