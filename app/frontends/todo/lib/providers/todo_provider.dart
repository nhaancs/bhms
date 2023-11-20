import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:todo/models/todo.dart';

final todoProvider = StateNotifierProvider<TodoListNotifier, List<Todo>>((ref) {
  return TodoListNotifier();
});

class TodoListNotifier extends StateNotifier<List<Todo>> {
  TodoListNotifier() : super([]);
  addTodo(String content) {
    state = [
      ...state,
      Todo(
        todoId: state.isEmpty ? 0 : state[state.length - 1].todoId + 1,
        content: content,
        completed: false,
      ),
    ];
  }

  completeTodo(int id) {
    state = [
      for (final todo in state)
        if (todo.todoId == id)
          Todo(
            todoId: todo.todoId,
            content: todo.content,
            completed: true,
          )
        else
          todo
    ];
  }

  deleteTodo(int id) {
    state = state.where((todo) => todo.todoId != id).toList();
  }
}
