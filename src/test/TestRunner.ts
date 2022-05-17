import Logger from "./lib/Logger";
import { AuthTestSuite } from "./auth/AuthTestSuite";

const TEST_SUITES = [
    AuthTestSuite
];


TEST_SUITES.forEach(testSuite => {
    new testSuite();
});

console.log = Logger.log;