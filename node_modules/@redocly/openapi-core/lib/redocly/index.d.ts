export declare class RedoclyClient {
    private accessToken;
    constructor();
    hasToken(): boolean;
    loadToken(): void;
    isAuthorizedWithRedocly(): Promise<boolean>;
    verifyToken(accessToken: string, verbose?: boolean): Promise<boolean>;
    getAuthorizationHeader(): Promise<string | undefined>;
    login(accessToken: string, verbose?: boolean): Promise<void>;
    logout(): void;
    query(queryString: string, parameters?: {}, headers?: {}): Promise<any>;
    static authorize(accessToken: string, options: {
        queryName?: string;
        verbose?: boolean;
    }): Promise<any>;
    updateDependencies(dependencies: string[] | undefined): Promise<void>;
    updateDefinitionVersion(definitionId: number, versionId: number, updatePatch: object): Promise<void>;
    getOrganizationId(organizationId: string): Promise<any>;
    getDefinitionByName(name: string, organizationId: string): Promise<any>;
    createDefinition(organizationId: string, name: string): Promise<any>;
    createDefinitionVersion(definitionId: string, name: string, sourceType: string, source: any): Promise<any>;
    getSignedUrl(organizationId: string, filesHash: string, fileName: string): Promise<any>;
    getDefinitionVersion(organizationId: string, definitionName: string, versionName: string): Promise<any>;
    static isRegistryURL(link: string): boolean;
}
