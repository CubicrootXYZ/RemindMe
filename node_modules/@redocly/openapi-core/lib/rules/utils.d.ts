import { UserContext } from '../walk';
export declare function oasTypeOf(value: unknown): "string" | "number" | "bigint" | "boolean" | "symbol" | "undefined" | "object" | "function" | "array" | "null";
/**
 * Checks if value matches specified JSON schema type
 *
 * @param {*} value - value to check
 * @param {JSONSchemaType} type - JSON Schema type
 * @returns boolean
 */
export declare function matchesJsonSchemaType(value: unknown, type: string, nullable: boolean): boolean;
export declare function missingRequiredField(type: string, field: string): string;
export declare function fieldNonEmpty(type: string, field: string): string;
export declare function validateDefinedAndNonEmpty(fieldName: string, value: any, ctx: UserContext): void;
export declare function getSuggest(given: string, variants: string[]): string[];
