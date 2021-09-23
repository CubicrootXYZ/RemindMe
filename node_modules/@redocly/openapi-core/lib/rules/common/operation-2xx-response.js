"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Operation2xxResponse = void 0;
const Operation2xxResponse = () => {
    return {
        ResponsesMap(responses, { report }) {
            const codes = Object.keys(responses);
            if (!codes.some((code) => code === 'default' || /2[Xx0-9]{2}/.test(code))) {
                report({
                    message: 'Operation must have at least one `2xx` response.',
                    location: { reportOnKey: true },
                });
            }
        },
    };
};
exports.Operation2xxResponse = Operation2xxResponse;
