import { ErrorObject } from '@redocly/ajv';
import { Location } from '../ref-utils';
import { ResolveFn } from '../walk';
export declare function releaseAjvInstance(): void;
export declare function validateJsonSchema(data: any, schema: any, schemaLoc: Location, instancePath: string, resolve: ResolveFn<any>, disallowAdditionalProperties: boolean): {
    valid: boolean;
    errors: (ErrorObject & {
        suggest?: string[];
    })[];
};
