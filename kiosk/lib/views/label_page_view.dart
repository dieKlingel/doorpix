import 'package:audioplayers/audioplayers.dart';
import 'package:flutter/material.dart';

import '../components/label_view.dart';
import '../core.dart';
import '../data/label_data.dart';

class LabelPageView extends StatelessWidget {
  final Label label;
  final Core core;
  final AudioPlayer player = AudioPlayer();
  final Source wav = AssetSource("bell.wav");

  LabelPageView({
    super.key,
    required this.core,
    required this.label,
  });

  void onPressed() {
    player.stop().then((value) => player.play(wav));

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
