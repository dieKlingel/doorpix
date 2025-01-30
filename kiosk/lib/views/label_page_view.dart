import 'package:flutter/material.dart';

import '../components/label_view.dart';
import '../core.dart';
import '../data/label_data.dart';

class LabelPageView extends StatelessWidget {
  final Label label;
  final Core core;

  const LabelPageView({
    super.key,
    required this.core,
    required this.label,
  });

  void onPressed() {
    core.emit("ring", data: {
      "source": "kiosk",
      "name": label.name,
      "street_name": label.street.name,
      "street_number": label.street.number,
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: LabelView(
        label,
        onPressed: onPressed,
      ),
    );
  }
}
