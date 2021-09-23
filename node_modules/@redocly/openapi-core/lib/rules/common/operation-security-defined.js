"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.OperationSecurityDefined = void 0;
const OperationSecurityDefined = () => {
    let referencedSchemes = new Map();
    return {
        DefinitionRoot: {
            leave(_, { report }) {
                for (const [name, scheme] of referencedSchemes.entries()) {
                    if (scheme.defined)
                        continue;
                    for (const reportedFromLocation of scheme.from) {
                        report({
                            message: `There is no \`${name}\` security scheme defined.`,
                            location: reportedFromLocation.key(),
                        });
                    }
                }
            },
        },
        SecurityScheme(_securityScheme, { key }) {
            referencedSchemes.set(key.toString(), {
                defined: true,
                from: [],
            });
        },
        SecurityRequirement(requirements, { location }) {
            for (const requirement of Object.keys(requirements)) {
                const authScheme = referencedSchemes.get(requirement);
                const requirementLocation = location.child([requirement]);
                if (!authScheme) {
                    referencedSchemes.set(requirement, { from: [requirementLocation] });
                }
                else {
                    authScheme.from.push(requirementLocation);
                }
            }
        },
    };
};
exports.OperationSecurityDefined = OperationSecurityDefined;
