import 'package:flutter/material.dart';

class QuestionScreen extends StatefulWidget {
  const QuestionScreen({super.key});

  @override
  State<QuestionScreen> createState() {
    return _QuestionScreenState();
  }
}

class _QuestionScreenState extends State<QuestionScreen> {
  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Text("the question"),
          const SizedBox(height: 30),
          ElevatedButton(onPressed: () {}, child: const Text('answer 1')),
          ElevatedButton(onPressed: () {}, child: const Text('answer 2')),
          ElevatedButton(onPressed: () {}, child: const Text('answer 3')),
          ElevatedButton(onPressed: () {}, child: const Text('answer 4')),
        ],
      ),
    );
  }
}
