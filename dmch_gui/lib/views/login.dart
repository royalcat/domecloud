import 'package:dmch_gui/bloc/auth.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

class LoginView extends StatefulWidget {
  const LoginView({Key? key}) : super(key: key);

  @override
  State<LoginView> createState() => _LoginViewState();
}

class _LoginViewState extends State<LoginView> {
  final _loginController = TextEditingController();
  final _passwordController = TextEditingController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Row(
        mainAxisSize: MainAxisSize.max,
        children: [
          Expanded(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Row(
                  children: const [
                    Spacer(),
                    Flexible(
                      flex: 3,
                      child: Image(
                        image: AssetImage('assets/images/cloud-icon.png'),
                      ),
                    ),
                    Spacer(),
                  ],
                ),
                const Text("DOMECLOUD"),
              ],
            ),
          ),
          const VerticalDivider(),
          Expanded(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                SizedBox(
                  width: 300,
                  child: Column(
                    children: [
                      TextField(
                        controller: _loginController,
                        decoration: const InputDecoration(
                          border: OutlineInputBorder(),
                          labelText: 'Login',
                        ),
                      ),
                      const SizedBox(height: 10),
                      TextField(
                        controller: _passwordController,
                        decoration: const InputDecoration(
                          border: OutlineInputBorder(),
                          labelText: 'Password',
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 30),
                ElevatedButton(
                  onPressed: () {
                    context.read<AuthenticationBloc>().add(
                          AuthenticationLogInEvent(
                            _loginController.text,
                            _passwordController.text,
                          ),
                        );
                  },
                  child: const SizedBox(
                    width: 100,
                    child: Center(
                      child: Text("Sign In"),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
