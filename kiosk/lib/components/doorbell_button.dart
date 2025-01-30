import 'package:flutter/material.dart';

class DoorbellButton extends StatelessWidget {
  final String text;
  final VoidCallback? onPressed;

  const DoorbellButton(
    this.text, {
    super.key,
    this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: ElevatedButton(
        onPressed: onPressed,
        style: ButtonStyle(
          side: WidgetStatePropertyAll(
            BorderSide(
              color: Colors.green,
              width: 3,
            ),
          ),
        ),
        child: ConstrainedBox(
          constraints: BoxConstraints(
            maxWidth: 460,
            minHeight: 100,
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              SizedBox(
                width: 80,
                child: Image.asset('assets/bell.png'),
              ),
              Expanded(
                child: Padding(
                  padding: const EdgeInsets.only(
                    right: 20,
                    top: 20,
                    bottom: 20,
                  ),
                  child: Text(
                    text,
                    textAlign: TextAlign.center,
                    textScaler: TextScaler.linear(1.2),
                  ),
                ),
              )
            ],
          ),
        ),

        //child: const Text('Ring Doorbell'),
      ),
    );
  }
}
