import { BaseResolver, Document } from './resolve';
import { NodeType } from './types';
import { LintConfig, Config } from './config/config';
export declare function lint(opts: {
    ref: string;
    config: Config;
    externalRefResolver?: BaseResolver;
}): Promise<import("./walk").NormalizedProblem[]>;
export declare function lintFromString(opts: {
    source: string;
    absoluteRef?: string;
    config: Config;
    externalRefResolver?: BaseResolver;
}): Promise<import("./walk").NormalizedProblem[]>;
export declare function lintDocument(opts: {
    document: Document;
    config: LintConfig;
    customTypes?: Record<string, NodeType>;
    externalRefResolver: BaseResolver;
}): Promise<import("./walk").NormalizedProblem[]>;
export declare function lintConfig(opts: {
    document: Document;
}): Promise<import("./walk").NormalizedProblem[]>;
