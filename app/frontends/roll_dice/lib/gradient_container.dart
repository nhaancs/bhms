import 'package:flutter/material.dart';

const startAlignment = Alignment.topLeft;
const endAlignment = Alignment.bottomRight;

class GradientContainer extends StatelessWidget {
  const GradientContainer({super.key, required this.colors});

  const GradientContainer.purple({super.key})
      : colors = const [Colors.deepPurple, Colors.indigo];

  final List<Color> colors;

  rollDice() {}

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: colors,
          begin: startAlignment,
          end: endAlignment,
        ),
      ),
      child: Center(
        child: Column(
          children: [
            Image.asset(
              "assets/images/dice-1.png",
              width: 200,
            ),
            TextButton(
              onPressed: rollDice,
              child: const Text("Roll Dice"),
            )
          ],
        ),
      ),
    );
  }
}
