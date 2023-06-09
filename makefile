
cleanup = docker-compose -p 'dissy' down > /dev/null;
# interruptClean should be used by tests that should clean up on interruption.
interruptClean = \
	echo "Shutting down gracefully - cleaning up"; \
	$(cleanup)

SetupD: Cleanup Build SetupDaux

SetupDaux:
	@trap '$(interruptClean)' INT TERM; \
	cd setupbooter; \
	go run main.go SetupD

Cleanup:
	@$(cleanup)

# > /dev/null 2>&1 is supressing output
Build:
	@docker-compose build > /dev/null 2>&1

