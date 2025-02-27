import {randomBytes} from 'crypto'
type TItem = {
    title: string;
    description: string;
    completed: boolean;
}
type TTodo = {
    //here we specify that once the id is set it can not be modified in the existing future
    readonly id?: string;
    readonly createdAt?: Date;
    item: TItem; 
}

//Here I create a typed string option that can be used to specify the expected errors during development
type TError = "createError" | "fetchError" | "updateError" | "deleteError" 

//Pick => used to pick attributes from existing types and return the type value of specified attribute
//Omit => used to omit certain attributes from a type and return a type without the specified type
interface Todo {
    addTodo(todo: TTodo): TTodo | TError
    getTodo(id: Pick<TTodo, "id"> | string): TTodo | TError
    updateTodo(id: Pick<TTodo, "id"> | string, data: Omit<TTodo, "id">): TTodo | TError
    deleteTodo(id: Pick<TTodo, "id"> | string): TError | void
}


class TodoGenerate implements Todo{
    private ids:string[] = [];

    public todos: TTodo[] = []
    constructor(){
        this.ids = [];
    }

    addTodo(todo: TTodo): TTodo | TError {
        const id = randomBytes(5).toString('hex')
        this.ids.push(id)
        const createdDate = new Date()
        const createdTodo: TTodo = {
            id: id,
            createdAt: createdDate,
            item: todo.item
        }
        this.todos.push(createdTodo);


        //simulate there's an error while adding a record
        if(!createdTodo){
            //this is automatically popped to you as a developer to ease development process
            return "createError"
        }

        return createdTodo;
    }

    getTodo(id: Pick<TTodo, 'id'> | string): TTodo | TError {
       const fetchedTodo = this.todos.find((todo) => id === todo.id) 

       if(!fetchedTodo){
            return "fetchError"
       }

       return fetchedTodo

    }

    updateTodo(id: Pick<TTodo, 'id'> | string, data: Partial<Omit<Partial<TTodo>, 'id'>>): TTodo | TError {
        
        let todo = this.getTodo(id);

        if(todo == "fetchError"){
            return "updateError"
        }

        todo = {
            ...data,
            ...todo as TTodo,
        }

        return todo;
    }   

    deleteTodo(id: Pick<TTodo, 'id'> | string): TError | void {
       let todo = this.getTodo(id);
       
       if(todo == "fetchError"){
            return "deleteError"
       }

       this.todos = this.todos.filter((todo) => todo.id !== id);

    }

}

const todoGenerator1 = new TodoGenerate();

let todos: TTodo[] = [
    {
        item: {
            title: "Wash clothes",
            description: "I am washing clothes today",
            completed: false
        }
    },
    {
        item: {
            title: "Play chess",
            description: "I am playing chess today",
            completed: false
        }
    },
    {
        item: {
            title: "Completing milestone 3 of my project",
            description: "I am completing milestone 3 of my  project",
            completed: false
        }
    },
    {
        item: {
            title: "Watch jujutsu kaisen as a reward of completing task 3", 
            description: "Rewarding yourself",
            completed: false
        }
    }


]

const todo1 = todoGenerator1.addTodo(
    {
        item: {
            title: "Wash clothes",
            description: "I am washing clothes today",
            completed: false
        }
    }
)
const todo2 = todoGenerator1.addTodo(
    {
        item: 
        {
            title: "Wash clothes",
            description: "I am washing clothes today",
            completed: false
        }
        
    }

) as TTodo

const todo2Updated = todoGenerator1.updateTodo(todo2.id as string, {
    item: {
        completed: true,
        title: '',
        description: ''
    }
})

 

