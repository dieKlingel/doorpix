import 'package:flutter/material.dart';
import 'package:kiosk/components/doorbell_button.dart';
import 'package:kiosk/components/street_view.dart';
import 'package:kiosk/data/label_data.dart';

class LabelView extends StatelessWidget {
  final Label data;
  final VoidCallback? onPressed;

  const LabelView(
    this.data, {
    super.key,
    this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      child: Card(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            StreetView(data.street),
            Padding(
              padding: const EdgeInsets.all(50.0),
              child: DoorbellButton(
                data.name,
                onPressed: onPressed,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
