import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_slidable/flutter_slidable.dart';
import 'package:todo/models/todo.dart';
import 'package:todo/pages/add.dart';
import 'package:todo/providers/todo_provider.dart';

class MyHomePage extends ConsumerWidget {
  const MyHomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    List<Todo> todos = ref.watch(todoProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text("Todo App"),
      ),
      body: ListView.builder(
        itemCount: todos.length,
        itemBuilder: (context, index) {
          return Slidable(
            startActionPane: ActionPane(
              motion: const ScrollMotion(),
              children: [
                SlidableAction(
                  onPressed: (context) => ref
                      .watch(todoProvider.notifier)
                      .deleteTodo(todos[index].todoId),
                  icon: Icons.delete,
                  backgroundColor: Colors.red,
                  borderRadius: const BorderRadius.all(Radius.circular(20)),
                )
              ],
            ),
            endActionPane: ActionPane(
              motion: const ScrollMotion(),
              children: [
                SlidableAction(
                  onPressed: (context) => ref
                      .watch(todoProvider.notifier)
                      .completeTodo(todos[index].todoId),
                  icon: Icons.check,
                  backgroundColor: Colors.green,
                  borderRadius: const BorderRadius.all(Radius.circular(20)),
                )
              ],
            ),
            child: ListTile(title: Text(todos[index].content)),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.of(context)
              .push(MaterialPageRoute(builder: (context) => const AddTodo()));
        },
        tooltip: 'Increment',
        child: const Icon(Icons.add),
      ), // This trailing comma makes auto-formatting nicer for build methods.
    );
  }
}
