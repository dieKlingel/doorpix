import 'package:flutter/material.dart';

import '../components/label_view.dart';
import '../data/label_data.dart';

class LabelPageView extends StatelessWidget {
  final Label label;

  const LabelPageView({
    super.key,
    required this.label,
  });

  void onPressed() {
    print('Doorbell button pressed');
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
