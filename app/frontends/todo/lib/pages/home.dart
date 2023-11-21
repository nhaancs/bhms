import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_slidable/flutter_slidable.dart';
import 'package:todo/models/todo.dart';
import 'package:todo/pages/add.dart';
import 'package:todo/pages/completed.dart';
import 'package:todo/providers/todo_provider.dart';

class MyHomePage extends ConsumerWidget {
  const MyHomePage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    List<Todo> todos = ref.watch(todoProvider);
    List<Todo> activeTodos = todos.where((todo) => !todo.completed).toList();
    List<Todo> completedTodos = todos.where((todo) => todo.completed).toList();

    return Scaffold(
      appBar: AppBar(
        title: const Text("Todo App"),
      ),
      body: ListView.builder(
        itemCount: activeTodos.length + 1,
        itemBuilder: (context, index) {
          if (index == activeTodos.length) {
            if (completedTodos.isEmpty) {
              return Container();
            }
            return Center(
              child: TextButton(
                child: const Text("Completed todos"),
                onPressed: () => Navigator.of(context).push(
                  MaterialPageRoute(
                      builder: (context) => const CompletedTodo()),
                ),
              ),
            );
          }
          return Slidable(
            startActionPane: ActionPane(
              motion: const ScrollMotion(),
              children: [
                SlidableAction(
                  onPressed: (context) => ref
                      .watch(todoProvider.notifier)
                      .deleteTodo(activeTodos[index].todoId),
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
                      .completeTodo(activeTodos[index].todoId),
                  icon: Icons.check,
                  backgroundColor: Colors.green,
                  borderRadius: const BorderRadius.all(Radius.circular(20)),
                )
              ],
            ),
            child: ListTile(title: Text(activeTodos[index].content)),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.of(context)
              .push(MaterialPageRoute(builder: (context) => const AddTodo()));
        },
        child: const Icon(Icons.add),
      ), // This trailing comma makes auto-formatting nicer for build methods.
    );
  }
}
