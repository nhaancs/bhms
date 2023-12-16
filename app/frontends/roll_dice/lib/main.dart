import 'package:flutter/material.dart';
import 'package:roll_dice/gradient_container.dart';

void main() {
  runApp(
    const MaterialApp(
      home: Scaffold(
        body: GradientContainer(
          colors: [
            Color.fromARGB(255, 80, 70, 2),
            Color.fromARGB(255, 2, 80, 68),
          ],
        ),
      ),
    ),
  );
}
