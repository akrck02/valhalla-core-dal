import Logger from "../lib/Logger";
import { TestState } from "./TestState";

export default class Test {

    private name : string;
    private runnable : Function;
    private state : TestState;

    constructor(runnable : Function) {
        this.name = runnable.name;
        this.runnable = runnable;
        this.state = TestState.PENDING;
    }

    public async run() : Promise<boolean> {
        
        const start = new Date().getTime();

        Logger.softTitle(this.name);

        try {
            await this.runnable();
            this.state = TestState.PASSED;

            const end = new Date().getTime();
            
            Logger.success("Test passed in " + (end-start) + "ms.");
            return true;
        } catch (e : any) {

            if(e.name === "AssertionError") {
                this.state = TestState.FAILED;
                Logger.warning("Test failed: " + this.name + " : " + e.message + "\n");
                process.exit(2);
                return false;
            }

            this.state = TestState.ERROR;
            Logger.error("Test terminated with an unexpected error");
            Logger.log(e);
            Logger.jump();
            process.exit(2);
            
            return false;
        }
    }

    public getState() : TestState {
        return this.state;
    }

    public getName() : string {
        return this.name;
    }

}