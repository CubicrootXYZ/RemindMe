"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.InfoDescription = void 0;
const utils_1 = require("../utils");
const InfoDescription = () => {
    return {
        Info(info, ctx) {
            utils_1.validateDefinedAndNonEmpty('description', info, ctx);
        },
    };
};
exports.InfoDescription = InfoDescription;
