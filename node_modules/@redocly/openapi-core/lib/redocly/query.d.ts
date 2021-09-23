export declare function query(query: string, variables?: {}, headers?: {}): Promise<any>;
export declare class GqlRequestError extends Error {
    constructor(message: string);
}
