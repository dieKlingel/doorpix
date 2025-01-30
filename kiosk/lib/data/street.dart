import 'package:yaml/yaml.dart';

class Street {
  final String name;
  final String number;

  const Street({
    required this.name,
    required this.number,
  });

  factory Street.fromYaml(YamlNode node) {
    if (node is! YamlMap) {
      throw ArgumentError.value(
        node,
        'the node is not a map',
      );
    }

    if (!node.containsKey('name')) {
      throw ArgumentError.value(
        node,
        'the node does not contain the key "name"',
      );
    }

    if (node['name'] is! String) {
      throw ArgumentError.value(
        node,
        'the value of the key "name" is not a string',
      );
    }

    if (!node.containsKey('number')) {
      throw ArgumentError.value(
        node,
        'the node does not contain the key number',
      );
    }

    if (node['number'] is! String) {
      throw ArgumentError.value(
        node,
        'the value of the key "number" is not a string',
      );
    }

    final name = node['name'] as String;
    final number = node['number'] as String;

    return Street(
      name: name.toString(),
      number: number.toString(),
    );
  }
}
