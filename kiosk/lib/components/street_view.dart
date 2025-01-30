import 'package:flutter/material.dart';

import '../data/street.dart';

class StreetView extends StatelessWidget {
  final Street street;

  const StreetView(
    this.street, {
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Text(
          street.number,
          textScaler: TextScaler.linear(8),
        ),
        Text(street.name),
      ],
    );
  }
}
