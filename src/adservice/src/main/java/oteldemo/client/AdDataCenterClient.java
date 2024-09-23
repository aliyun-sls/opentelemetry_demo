package oteldemo.client;

import okhttp3.OkHttpClient;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import retrofit2.Call;
import retrofit2.Response;
import retrofit2.Retrofit;
import retrofit2.converter.gson.GsonConverterFactory;
import retrofit2.converter.scalars.ScalarsConverterFactory;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;

public class AdDataCenterClient {

    private static final Logger logger = LogManager.getLogger(AdDataCenterClient.class);

    private static final IComponentAPI client = initClient();

    private static final AdDataCenterClient instance = new AdDataCenterClient();

    private static IComponentAPI initClient() {
        int timeOut = 5;
        OkHttpClient okHttpClient = new OkHttpClient.Builder()
                .readTimeout(timeOut, TimeUnit.SECONDS)
                .connectTimeout(timeOut, TimeUnit.SECONDS)
                .writeTimeout(timeOut, TimeUnit.SECONDS)
                .build();

        String endpoint = System.getenv("AD_DATA_CENTER_ENDPOINT");

        Retrofit retrofit = new Retrofit.Builder()
                .baseUrl(endpoint)
                .client(okHttpClient)
                .addConverterFactory(ScalarsConverterFactory.create())
                .addConverterFactory(GsonConverterFactory.create())
                .build();

        return retrofit.create(IComponentAPI.class);
    }

    public static AdDataCenterClient getInstance() {
        return instance;
    }

    public void getData() {
        Invocation invocation = buildInvocation();
        List<Invocation> invocations = new ArrayList<>();
        invocations.add(invocation);
        Call<String> call = client.invokeChildren("ads", "data", invocations, 0);
        try {
            Response<String> response = call.execute();
            logger.info("execute mock server call success, response: {}", response.body());
        } catch (Throwable t) {
            logger.error("execute mock server call error", t);
        }
    }

    private Invocation buildInvocation() {
//        Invocation datacenter = new Invocation("ad-datacenter", "ads", "data", new ArrayList<>());
        Invocation cacheService = new Invocation("ad-cache-service", "ads", "cache", new ArrayList<>());
        Invocation dbService = new Invocation("ad-db-service", "ads", "db", new ArrayList<>());
//        datacenter.getChildren().add(cacheService);
        cacheService.getChildren().add(dbService);
        return cacheService;
    }
}
