import Logger from "../lib/Logger";
import Test from "./Test";

export default abstract class TestSuite {

    constructor(tests : Test[]){
        TestSuite.runAll(this.constructor.name,tests);
    }

    static async runAll(name: string, tests : Test[]) {

        const start = new Date().getTime();

        Logger.hardTitle("🧪 Running test suite " + name + " 🧪");

        for(let i = 0; i < tests.length; i++) {
            
            const testResult = await tests[i].run();
            
            if(!testResult) {
                Logger.softTitle("🟡 Tests terminated with a failure 🟡");
                return;
            }
        }

        const end = new Date().getTime();

        Logger.softTitle("🟢 All tests passed successfully 🟢");
        Logger.rawlog("Test suite ran in " + (end-start) + "ms.\n");
        
    }

} 