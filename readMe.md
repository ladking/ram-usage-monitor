### RAM USAGE MONITOR ###

Simple program to monitor the RAM usage on my Linux (Ubuntu) computer. The program runs at set interval, checks if the current RAM usage exceeds the set limit, then sends a desktop alert so i can close processes consuming alot of memory to prevent the system from freezing.

Might extend the code to automatically find processes consuming memory and kill them, but the current logic serves the purpose i need it for.

The program runs as a daemon managed by systemd on my computer