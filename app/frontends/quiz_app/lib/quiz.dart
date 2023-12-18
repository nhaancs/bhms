import 'package:flutter/material.dart';
import 'package:quiz_app/data/questions.dart';
import 'package:quiz_app/question_screen.dart';
import 'package:quiz_app/results_screen.dart';
import 'package:quiz_app/start_screen.dart';

const startScreen = "start-screen";
const questionScreen = "question-screen";
const resultsScreen = "results-screen";

class Quiz extends StatefulWidget {
  const Quiz({super.key});

  @override
  State<Quiz> createState() {
    return _QuizState();
  }
}

class _QuizState extends State<Quiz> {
  var activeScreen = startScreen;
  // private property
  final List<String> _selectedAnswers = [];

  void chooseAnswer(String answer) {
    _selectedAnswers.add(answer);

    if (_selectedAnswers.length >= questions.length) {
      setState(() {
        activeScreen = resultsScreen;
      });
    }
  }

  void switchToQuestionScreen() {
    setState(() {
      activeScreen = questionScreen;
    });
  }

  @override
  Widget build(BuildContext context) {
    Widget? screen;
    switch (activeScreen) {
      case startScreen:
        screen = StartScreen(switchToQuestionScreen);
      case questionScreen:
        screen = QuestionScreen(onSelectAnswer: chooseAnswer);
      case resultsScreen:
        screen = ResultsScreen(chosenAnswers: _selectedAnswers);
      default:
        screen = StartScreen(switchToQuestionScreen);
    }

    return MaterialApp(
      home: Scaffold(
        body: Container(
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              colors: [
                Color.fromARGB(255, 78, 13, 151),
                Color.fromARGB(255, 107, 15, 168),
              ],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
          ),
          child: screen,
        ),
      ),
    );
  }
}
