"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NoIdenticalPaths = void 0;
const NoIdenticalPaths = () => {
    return {
        PathMap(pathMap, { report, location }) {
            const pathsMap = new Map();
            for (const pathName of Object.keys(pathMap)) {
                const id = pathName.replace(/{.+?}/g, '{VARIABLE}');
                const existingSamePath = pathsMap.get(id);
                if (existingSamePath) {
                    report({
                        message: `The path already exists which differs only by path parameter name(s): \`${existingSamePath}\` and \`${pathName}\`.`,
                        location: location.child([pathName]).key(),
                    });
                }
                else {
                    pathsMap.set(id, pathName);
                }
            }
        },
    };
};
exports.NoIdenticalPaths = NoIdenticalPaths;
