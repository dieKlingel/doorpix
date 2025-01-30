import 'package:kiosk/data/street.dart';
import 'package:yaml/yaml.dart';

class Label {
  final String name;
  final Street street;

  const Label({
    required this.name,
    required this.street,
  });

  factory Label.fromYaml(YamlNode node) {
    if (node is! YamlMap) {
      throw ArgumentError.value(
        node,
        'the node is not a map',
      );
    }

    if (!node.containsKey('title')) {
      throw ArgumentError.value(
        node,
        'the node does not contain the key "title"',
      );
    }

    if (node['title'] is! String) {
      throw ArgumentError.value(
        node,
        'the value of the key "title" is not a string',
      );
    }

    if (!node.containsKey('street')) {
      throw ArgumentError.value(
        node,
        'the node does not contain the key "street"',
      );
    }

    if (node['street'] is! YamlMap) {
      throw ArgumentError.value(
        node,
        'the value of the key street is not a map',
      );
    }

    final name = node['title'] as String;
    final street = Street.fromYaml(node['street']);

    return Label(
      name: name,
      street: street,
    );
  }
}
