import 'package:flutter/material.dart';
import 'package:roll_dice/gradient_container.dart';

void main() {
  runApp(
    MaterialApp(
      home: Scaffold(
        body: GradientContainer(
          colors: const [
            Color.fromARGB(255, 80, 70, 2),
            Color.fromARGB(255, 2, 80, 68),
          ],
        ),
      ),
    ),
  );
}
