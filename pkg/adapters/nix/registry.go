package nix

import (
	"os/exec"

	"go.uber.org/zap"
)

func (nix *Nix) AddRegistry(alias, registry string) error {
	err := nix.addRegistry(alias, registry)
	if err != nil {
		return err
	}

	err = nix.validateRegistry(alias)
	if err != nil {
		return nix.RemoveRegistry(alias)
	}

	return nix.PinRegistry(registry)
}

func (nix *Nix) RemoveRegistry(alias string) error {
	args := append(
		nix.flags,
		"registry",
		"remove",
		alias,
	)

	cmd := exec.Command(nix.command, args...) //nolint
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		err := reflectNixError(string(cmdOutput))
		if err == ErrUnknownNixError {
			nix.logger.Error(
				"failed to execute nix registry remove",
				zap.Error(err), zap.Strings("args", args), zap.ByteString("cmdOutput", cmdOutput),
			)
		}
		return err
	}

	return nil
}

func (nix *Nix) PinRegistry(registry string) error {
	args := append(
		nix.flags,
		"registry",
		"pin",
		registry,
	)

	cmd := exec.Command(nix.command, args...) //nolint
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		err := reflectNixError(string(cmdOutput))
		if err == ErrUnknownNixError {
			nix.logger.Error(
				"failed to execute nix registry pin",
				zap.Error(err), zap.Strings("args", args), zap.ByteString("cmdOutput", cmdOutput),
			)
		}
		return err
	}

	return nil
}

func (nix *Nix) addRegistry(alias, registry string) error {
	args := append(
		nix.flags,
		"registry",
		"add",
		alias,
		registry,
	)

	cmd := exec.Command(nix.command, args...) //nolint
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		err := reflectNixError(string(cmdOutput))
		if err == ErrUnknownNixError {
			nix.logger.Error(
				"failed to execute nix registry add",
				zap.Error(err), zap.Strings("args", args), zap.ByteString("cmdOutput", cmdOutput),
			)
		}
		return err
	}
	return nil
}

func (nix *Nix) validateRegistry(alias string) error {
	args := append(
		nix.flags,
		"flake",
		"metadata",
		alias,
	)

	cmd := exec.Command(nix.command, args...) //nolint
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		err := reflectNixError(string(cmdOutput))
		if err == ErrUnknownNixError {
			nix.logger.Error(
				"failed to execute nix flake metadata",
				zap.Error(err), zap.Strings("args", args), zap.ByteString("cmdOutput", cmdOutput),
			)
		}
		return err
	}

	return nil
}
