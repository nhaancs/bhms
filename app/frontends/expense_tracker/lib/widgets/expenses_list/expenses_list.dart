import 'package:expense_tracker/models/expense.dart';
import 'package:expense_tracker/widgets/expenses_list/expense_item.dart';
import 'package:flutter/material.dart';

class ExpensesList extends StatelessWidget {
  const ExpensesList({
    super.key,
    required this.expenses,
    required this.onRemoveExpense,
  });

  final List<Expense> expenses;
  final void Function(Expense) onRemoveExpense;

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: ListView.builder(
        itemCount: expenses.length,
        itemBuilder: (ctx, index) => Dismissible(
            background: Container(
              color: Theme.of(context).colorScheme.error.withOpacity(0.8),
              margin: EdgeInsets.symmetric(
                  horizontal: Theme.of(context).cardTheme.margin!.horizontal),
            ),
            key: ValueKey(expenses[index]),
            onDismissed: (_) {
              onRemoveExpense(expenses[index]);
            },
            child: ExpenseItem(expenses[index])),
      ),
    );
  }
}
