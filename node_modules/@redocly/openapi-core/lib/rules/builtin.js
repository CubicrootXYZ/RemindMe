"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.decorators = exports.preprocessors = exports.rules = void 0;
const oas3 = require("./oas3/index");
const oas2 = require("./oas2/index");
exports.rules = {
    oas3: oas3.rules,
    oas2: oas2.rules,
};
exports.preprocessors = {
    oas3: oas3.preprocessors,
    oas2: oas2.preprocessors,
};
exports.decorators = {
    oas3: oas3.decorators,
    oas2: oas2.decorators,
};
