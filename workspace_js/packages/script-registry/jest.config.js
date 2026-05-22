module.exports = {
  globals: {
    window: {}
  },
  moduleFileExtensions: [
    "js",
    "json",
    "ts"
  ],
  rootDir: "src",
  testRegex: ".*\\.spec\\.ts$",
  transform: {
    "^.+\\.(t|j)s$": ["ts-jest", {
      tsconfig: 'tsconfig.json'
    }]
  },
  collectCoverageFrom: [
    "**/*.(t|j)s"
  ],
  coverageDirectory: "../coverage",
  testEnvironment: "node",
  moduleNameMapper: {
    "^@/(.*)$": "<rootDir>/../src/$1",
    "^@src/(.*)$": "<rootDir>/../src/$1",
  }
}
