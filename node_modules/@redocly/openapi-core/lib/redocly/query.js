"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.GqlRequestError = exports.query = void 0;
const node_fetch_1 = require("node-fetch");
const GRAPHQL_ENDPOINT = process.env.REDOCLY_DOMAIN
    ? `https://api.${process.env.REDOCLY_DOMAIN}/graphql`
    : 'https://api.redoc.ly/graphql';
function query(query, variables = {}, headers = {}) {
    return __awaiter(this, void 0, void 0, function* () {
        headers = Object.assign(Object.assign({}, headers), { 'Content-Type': 'application/json' });
        const gQLResponse = yield node_fetch_1.default(GRAPHQL_ENDPOINT, {
            method: 'POST',
            headers,
            body: JSON.stringify({
                query,
                variables,
            }),
        });
        if (!gQLResponse.ok) {
            throw new GqlRequestError(`Failed to execute query: ${gQLResponse.status}`);
        }
        const response = yield gQLResponse.json();
        if (response.errors && response.errors.length) {
            throw new GqlRequestError(`Query failed: ${response.errors[0].message}`);
        }
        return response.data;
    });
}
exports.query = query;
class GqlRequestError extends Error {
    constructor(message) {
        super(message);
    }
}
exports.GqlRequestError = GqlRequestError;
