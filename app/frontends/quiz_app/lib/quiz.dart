import 'package:flutter/material.dart';
import 'package:quiz_app/data/questions.dart';
import 'package:quiz_app/question_screen.dart';
import 'package:quiz_app/results_screen.dart';
import 'package:quiz_app/start_screen.dart';

const startScreen = "start-screen";
const questionsScreen = "questions-screen";
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
  List<String> selectedAnswers = [];

  void chooseAnswer(String answer) {
    selectedAnswers.add(answer);

    if (selectedAnswers.length >= questions.length) {
      setState(() {
        selectedAnswers = [];
        activeScreen = resultsScreen;
      });
    }
  }

  void switchToQuestionScreen() {
    setState(() {
      activeScreen = questionsScreen;
    });
  }

  @override
  Widget build(BuildContext context) {
    Widget? screen;
    switch (activeScreen) {
      case startScreen:
        screen = StartScreen(switchToQuestionScreen);
      case questionsScreen:
        screen = QuestionScreen(onSelectAnswer: chooseAnswer);
      case resultsScreen:
        screen = const ResultsScreen();
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
