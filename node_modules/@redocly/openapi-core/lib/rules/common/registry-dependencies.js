"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.RegistryDependencies = void 0;
const redocly_1 = require("../../redocly");
const RegistryDependencies = () => {
    let redoclyClient;
    let registryDependencies = new Set();
    return {
        DefinitionRoot: {
            leave() {
                redoclyClient = new redocly_1.RedoclyClient();
                if (process.env.UPDATE_REGISTRY && redoclyClient.hasToken()) {
                    redoclyClient.updateDependencies(Array.from(registryDependencies.keys()));
                }
            },
        },
        ref(node) {
            if (node.$ref) {
                const link = node.$ref.split('#/')[0];
                if (redocly_1.RedoclyClient.isRegistryURL(link)) {
                    registryDependencies.add(link);
                }
            }
        },
    };
};
exports.RegistryDependencies = RegistryDependencies;
