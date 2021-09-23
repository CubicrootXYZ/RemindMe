import type { Oas3Definition, Oas3ExternalDocs, Oas3Info, Oas3Contact, Oas3Components, Oas3License, Oas3Schema, Oas3Header, Oas3Parameter, Oas3Operation, Oas3PathItem, Oas3ServerVariable, Oas3Server, Oas3MediaType, Oas3Response, Oas3Example, Oas3RequestBody, Oas3Tag, OasRef, Oas3SecurityScheme, Oas3SecurityRequirement, Oas3Encoding, Oas3Link, Oas3Xml, Oas3Discriminator, Oas3Callback } from './typings/openapi';
import { Oas2Definition, Oas2Tag, Oas2ExternalDocs, Oas2SecurityRequirement, Oas2Info, Oas2Contact, Oas2License, Oas2PathItem, Oas2Operation, Oas2Header, Oas2Response, Oas2Schema, Oas2Xml, Oas2Parameter, Oas2SecurityScheme } from './typings/swagger';
import { NormalizedNodeType } from './types';
import { Stack } from './utils';
import { UserContext, ResolveResult, ProblemSeverity } from './walk';
import { Location } from './ref-utils';
export declare type VisitFunction<T> = (node: T, ctx: UserContext & {
    ignoreNextVisitorsOnNode: () => void;
}, parents?: any, context?: any) => void;
declare type VisitRefFunction = (node: OasRef, ctx: UserContext, resolved: ResolveResult<any>) => void;
declare type SkipFunction<T> = (node: T, key: string | number) => boolean;
declare type VisitObject<T> = {
    enter?: VisitFunction<T>;
    leave?: VisitFunction<T>;
    skip?: SkipFunction<T>;
};
declare type NestedVisitObject<T, P> = VisitObject<T> & NestedVisitor<P>;
declare type VisitFunctionOrObject<T> = VisitFunction<T> | VisitObject<T>;
declare type VisitorNode<T extends any> = {
    ruleId: string;
    severity: ProblemSeverity;
    context: VisitorLevelContext | VisitorSkippedLevelContext;
    depth: number;
    visit: VisitFunction<T>;
    skip?: SkipFunction<T>;
};
declare type VisitorRefNode = {
    ruleId: string;
    severity: ProblemSeverity;
    context: VisitorLevelContext;
    depth: number;
    visit: VisitRefFunction;
};
export declare type VisitorLevelContext = {
    isSkippedLevel: false;
    type: NormalizedNodeType;
    parent: VisitorLevelContext | null;
    activatedOn: Stack<{
        node?: any;
        withParentNode?: any;
        skipped: boolean;
        nextLevelTypeActivated: Stack<NormalizedNodeType>;
        location?: Location;
    }>;
};
export declare type VisitorSkippedLevelContext = {
    isSkippedLevel: true;
    parent: VisitorLevelContext;
    seen: Set<any>;
};
declare type NormalizeVisitor<Fn> = Fn extends VisitFunction<infer T> ? VisitorNode<T> : never;
export declare type BaseVisitor = {
    any?: {
        enter?: VisitFunction<any>;
        leave?: VisitFunction<any>;
        skip?: SkipFunction<any>;
    } | VisitFunction<any>;
    ref?: {
        enter?: VisitRefFunction;
        leave?: VisitRefFunction;
    } | VisitRefFunction;
};
declare type Oas3FlatVisitor = {
    DefinitionRoot?: VisitFunctionOrObject<Oas3Definition>;
    Tag?: VisitFunctionOrObject<Oas3Tag>;
    ExternalDocs?: VisitFunctionOrObject<Oas3ExternalDocs>;
    Server?: VisitFunctionOrObject<Oas3Server>;
    ServerVariable?: VisitFunctionOrObject<Oas3ServerVariable>;
    SecurityRequirement?: VisitFunctionOrObject<Oas3SecurityRequirement>;
    Info?: VisitFunctionOrObject<Oas3Info>;
    Contact?: VisitFunctionOrObject<Oas3Contact>;
    License?: VisitFunctionOrObject<Oas3License>;
    PathMap?: VisitFunctionOrObject<Record<string, Oas3PathItem>>;
    PathItem?: VisitFunctionOrObject<Oas3PathItem>;
    Callback?: VisitFunctionOrObject<Record<string, Oas3PathItem>>;
    Parameter?: VisitFunctionOrObject<Oas3Parameter>;
    Operation?: VisitFunctionOrObject<Oas3Operation>;
    RequestBody?: VisitFunctionOrObject<Oas3RequestBody>;
    MediaTypeMap?: VisitFunctionOrObject<Record<string, Oas3MediaType>>;
    MediaType?: VisitFunctionOrObject<Oas3MediaType>;
    Example?: VisitFunctionOrObject<Oas3Example>;
    Encoding?: VisitFunctionOrObject<Oas3Encoding>;
    Header?: VisitFunctionOrObject<Oas3Header>;
    ResponsesMap?: VisitFunctionOrObject<Record<string, Oas3Response>>;
    Response?: VisitFunctionOrObject<Oas3Response>;
    Link?: VisitFunctionOrObject<Oas3Link>;
    Schema?: VisitFunctionOrObject<Oas3Schema>;
    Xml?: VisitFunctionOrObject<Oas3Xml>;
    SchemaProperties?: VisitFunctionOrObject<Record<string, Oas3Schema>>;
    DiscriminatorMapping?: VisitFunctionOrObject<Record<string, string>>;
    Discriminator?: VisitFunctionOrObject<Oas3Discriminator>;
    Components?: VisitFunctionOrObject<Oas3Components>;
    NamedSchemas?: VisitFunctionOrObject<Record<string, Oas3Schema>>;
    NamedResponses?: VisitFunctionOrObject<Record<string, Oas3Response>>;
    NamedParameters?: VisitFunctionOrObject<Record<string, Oas3Parameter>>;
    NamedExamples?: VisitFunctionOrObject<Record<string, Oas3Example>>;
    NamedRequestBodies?: VisitFunctionOrObject<Record<string, Oas3RequestBody>>;
    NamedHeaders?: VisitFunctionOrObject<Record<string, Oas3Header>>;
    NamedSecuritySchemes?: VisitFunctionOrObject<Record<string, Oas3SecurityScheme>>;
    NamedLinks?: VisitFunctionOrObject<Record<string, Oas3Link>>;
    NamedCallbacks?: VisitFunctionOrObject<Record<string, Oas3Callback>>;
    ImplicitFlow?: VisitFunctionOrObject<Oas3SecurityScheme['flows']['implicit']>;
    PasswordFlow?: VisitFunctionOrObject<Oas3SecurityScheme['flows']['password']>;
    ClientCredentials?: VisitFunctionOrObject<Oas3SecurityScheme['flows']['clientCredentials']>;
    AuthorizationCode?: VisitFunctionOrObject<Oas3SecurityScheme['flows']['authorizationCode']>;
    SecuritySchemeFlows?: VisitFunctionOrObject<Oas3SecurityScheme['flows']>;
    SecurityScheme?: VisitFunctionOrObject<Oas3SecurityScheme>;
};
declare type Oas2FlatVisitor = {
    DefinitionRoot?: VisitFunctionOrObject<Oas2Definition>;
    Tag?: VisitFunctionOrObject<Oas2Tag>;
    ExternalDocs?: VisitFunctionOrObject<Oas2ExternalDocs>;
    SecurityRequirement?: VisitFunctionOrObject<Oas2SecurityRequirement>;
    Info?: VisitFunctionOrObject<Oas2Info>;
    Contact?: VisitFunctionOrObject<Oas2Contact>;
    License?: VisitFunctionOrObject<Oas2License>;
    PathMap?: VisitFunctionOrObject<Record<string, Oas2PathItem>>;
    PathItem?: VisitFunctionOrObject<Oas2PathItem>;
    Parameter?: VisitFunctionOrObject<any>;
    Operation?: VisitFunctionOrObject<Oas2Operation>;
    Examples?: VisitFunctionOrObject<Record<string, any>>;
    Header?: VisitFunctionOrObject<Oas2Header>;
    ResponsesMap?: VisitFunctionOrObject<Record<string, Oas2Response>>;
    Response?: VisitFunctionOrObject<Oas2Response>;
    Schema?: VisitFunctionOrObject<Oas2Schema>;
    Xml?: VisitFunctionOrObject<Oas2Xml>;
    SchemaProperties?: VisitFunctionOrObject<Record<string, Oas2Schema>>;
    NamedSchemas?: VisitFunctionOrObject<Record<string, Oas2Schema>>;
    NamedResponses?: VisitFunctionOrObject<Record<string, Oas2Response>>;
    NamedParameters?: VisitFunctionOrObject<Record<string, Oas2Parameter>>;
    SecurityScheme?: VisitFunctionOrObject<Oas2SecurityScheme>;
};
declare type Oas3NestedVisitor = {
    [T in keyof Oas3FlatVisitor]: Oas3FlatVisitor[T] extends Function ? Oas3FlatVisitor[T] : Oas3FlatVisitor[T] & NestedVisitor<Oas3NestedVisitor>;
};
declare type Oas2NestedVisitor = {
    [T in keyof Oas2FlatVisitor]: Oas2FlatVisitor[T] extends Function ? Oas2FlatVisitor[T] : Oas2FlatVisitor[T] & NestedVisitor<Oas2NestedVisitor>;
};
export declare type Oas3Visitor = BaseVisitor & Oas3NestedVisitor & Record<string, VisitFunction<any> | NestedVisitObject<any, Oas3NestedVisitor>>;
export declare type Oas2Visitor = BaseVisitor & Oas2NestedVisitor & Record<string, VisitFunction<any> | NestedVisitObject<any, Oas2NestedVisitor>>;
export declare type Oas3TransformVisitor = BaseVisitor & Oas3FlatVisitor & Record<string, VisitFunction<any> | VisitObject<any>>;
export declare type Oas2TransformVisitor = BaseVisitor & Oas2FlatVisitor & Record<string, VisitFunction<any> | VisitObject<any>>;
export declare type NestedVisitor<T> = Exclude<T, 'any' | 'ref' | 'DefinitionRoot'>;
export declare type NormalizedOasVisitors<T extends BaseVisitor> = {
    [V in keyof T]-?: {
        enter: Array<NormalizeVisitor<T[V]>>;
        leave: Array<NormalizeVisitor<T[V]>>;
    };
} & {
    ref: {
        enter: Array<VisitorRefNode>;
        leave: Array<VisitorRefNode>;
    };
    [k: string]: {
        enter: Array<VisitorNode<any>>;
        leave: Array<VisitorNode<any>>;
    };
};
export declare type Oas3Rule = (options: Record<string, any>) => Oas3Visitor;
export declare type Oas2Rule = (options: Record<string, any>) => Oas2Visitor;
export declare type Oas3Preprocessor = (options: Record<string, any>) => Oas3TransformVisitor;
export declare type Oas2Preprocessor = (options: Record<string, any>) => Oas2TransformVisitor;
export declare type Oas3Decorator = (options: Record<string, any>) => Oas3TransformVisitor;
export declare type Oas2Decorator = (options: Record<string, any>) => Oas2TransformVisitor;
export declare type OasRule = Oas3Rule;
export declare type OasPreprocessor = Oas3Preprocessor;
export declare type OasDecorator = Oas3Decorator;
export declare type RuleInstanceConfig = {
    ruleId: string;
    severity: ProblemSeverity;
};
export declare function normalizeVisitors<T extends BaseVisitor>(visitorsConfig: (RuleInstanceConfig & {
    visitor: NestedVisitObject<any, T>;
})[], types: Record<keyof T, NormalizedNodeType>): NormalizedOasVisitors<T>;
export {};
