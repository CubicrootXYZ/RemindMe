"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NoEmptyServers = void 0;
const NoEmptyServers = () => {
    return {
        DefinitionRoot(root, { report, location }) {
            if (!root.servers) {
                report({
                    message: 'Servers must be present.',
                });
                return;
            }
            if (!Array.isArray(root.servers) || root.servers.length === 0) {
                report({
                    message: 'Servers must be a non-empty array.',
                    location: location.child(['servers']).key(),
                });
            }
        },
    };
};
exports.NoEmptyServers = NoEmptyServers;
