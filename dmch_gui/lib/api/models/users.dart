// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'dart:convert';

import 'package:equatable/equatable.dart';

class User extends Equatable {
  final String name;
  final bool isAdmin;
  final String password;

  const User({
    required this.name,
    required this.isAdmin,
    required this.password,
  });

  const User.undefined()
      : this(
          name: '',
          isAdmin: false,
          password: '',
        );

  User copyWith({
    String? name,
    bool? isAdmin,
    String? password,
  }) {
    return User(
      name: name ?? this.name,
      isAdmin: isAdmin ?? this.isAdmin,
      password: password ?? this.password,
    );
  }

  Map<String, dynamic> toMap() {
    return <String, dynamic>{
      'name': name,
      'isAdmin': isAdmin,
      'password': password,
    };
  }

  factory User.fromMap(Map<String, dynamic> map) {
    return User(
      name: map['username'] as String,
      isAdmin: map['isAdmin'] as bool,
      password: map['password'] as String,
    );
  }

  String toJson() => json.encode(toMap());

  factory User.fromJson(String source) => User.fromMap(json.decode(source) as Map<String, dynamic>);

  @override
  bool get stringify => true;

  @override
  List<Object> get props => [name, isAdmin, password];
}
