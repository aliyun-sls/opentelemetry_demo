package oteldemo.problempattern;

import org.apache.commons.lang3.StringUtils;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

import java.sql.*;

public enum DatabaseCaller {

    INSTANCE;

    private final Logger logger = LogManager.getLogger(DatabaseCaller.class);

    private final String url;
    private final String user;
    private final String password;

    DatabaseCaller() {
        try {
            Class.forName("com.mysql.cj.jdbc.Driver");
        } catch (ClassNotFoundException e) {
            throw new RuntimeException(e);
        }

        url = System.getenv("DB_URL");
        user = System.getenv("DB_USER");
        password = System.getenv("DB_PASSWORD");

        logger.info("DatabaseCaller initialized. {}, {}, {}", url, user, password);
    }

    public void execute(Boolean withNotExistTable) throws SQLException {
        if (StringUtils.isEmpty(url)) {
            logger.warn("DatabaseCaller is not initialized, skip to execute");
            return;
        }

        String query = String.format("SELECT * FROM %s", withNotExistTable ? "users_1" : "users");
        logger.info("Executing query: {}", query);
        try (Connection conn = DriverManager.getConnection(url, user, password); Statement stmt = conn.createStatement(); ResultSet rs = stmt.executeQuery(query)) {
            while (rs.next()) {
                int id = rs.getInt("id");
                String name = rs.getString("name");
                logger.info("id: {}, name: {}", id, name);
            }
        }
    }

}
