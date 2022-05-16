export default interface ITask {
    id: number;
    author: string;
    name: string;
    description?: string;
    start?: string;
    end?: string;
    allDay?: number;
    done?: number;
    labels?: string[];
}