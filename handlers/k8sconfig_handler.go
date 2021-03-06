package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/layer5io/meshery/meshes"
	"github.com/sirupsen/logrus"
)

func (h *Handler) K8SConfigHandler(w http.ResponseWriter, req *http.Request) {
	// 		if r.Method == http.MethodGet {
	// 			data := map[string]interface{}{
	// 				"ByPassAuth": h.config.ByPassAuth,
	// 			}

	// 			session, err := h.config.SessionStore.Get(r, h.config.SessionName)
	// 			if err != nil {
	// 				logrus.Errorf("error getting session: %v", err)
	// 				http.Error(w, "unable to get session", http.StatusUnauthorized)
	// 				return
	// 			}

	// 			if !h.config.ByPassAuth {
	// 				user, _ := session.Values["user"].(*models.User)
	// 				data["User"] = user
	// 			}

	// 			data["Flashes"] = session.Flashes()
	// 			session.Save(r, w)

	// 			err = getK8SConfigTempl.Execute(w, data)
	// 			if err != nil {
	// 				logrus.Errorf("error rendering the template for the page: %v", err)
	// 				http.Error(w, "unable to serve the requested file", http.StatusInternalServerError)
	// 				return
	// 			}
	// 		} else if r.Method == http.MethodPost {
	// 			h.DashboardHandler(ctx, w, r)
	// 		} else {
	// 			w.WriteHeader(http.StatusNotFound)
	// 		}
	// 	}
	// }

	// func (h *Handler) DashboardHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(1 << 20)

	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	// var user *models.User
	// if h.config.ByPassAuth {
	// 	userName := req.FormValue("user_name")
	// 	if userName == "" {
	// 		userName = "Test User"
	// 	}
	// 	user = h.setupSession(userName, req, w)
	// 	if user == nil {
	// 		return
	// 	}
	// } else {
	// 	// session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	// 	// if err != nil {
	// 	// 	logrus.Errorf("error getting session: %v", err)
	// 	// 	http.Error(w, "unable to get session", http.StatusUnauthorized)
	// 	// 	return
	// 	// }
	// 	user, _ = session.Values["user"].(*models.User)
	// }
	inClusterConfig := req.FormValue("inClusterConfig")
	logrus.Debugf("inClusterConfig: %s", inClusterConfig)

	var k8sConfigBytes []byte
	var contextName string

	if inClusterConfig == "" {
		// k8sfile, contextName
		k8sfile, _, err := req.FormFile("k8sfile")
		if err != nil {
			logrus.Errorf("error getting k8s file: %v", err)
			// http.Error(w, "error getting k8s file", http.StatusUnauthorized)
			// session.AddFlash("Unable to get kubernetes config file")
			// session.Save(req, w)
			// http.Redirect(w, req, "/play/dashboard", http.StatusFound)
			http.Error(w, "Unable to get kubernetes config file", http.StatusBadRequest)
			return
		}
		defer k8sfile.Close()
		k8sConfigBytes, err = ioutil.ReadAll(k8sfile)
		if err != nil {
			logrus.Errorf("error reading config: %v", err)
			// http.Error(w, "unable to read config", http.StatusBadRequest)
			// session.AddFlash("Unable to read the kubernetes config file, please try again")
			// session.Save(req, w)
			// http.Redirect(w, req, "/play/dashboard", http.StatusFound)
			http.Error(w, "Unable to read the kubernetes config file, please try again", http.StatusBadRequest)
			return
		}

		contextName = req.FormValue("contextName")

		ccfg, err := clientcmd.Load(k8sConfigBytes)
		if err != nil {
			logrus.Errorf("error parsing k8s config: %v", err)
			// http.Error(w, "k8s config file not valid", http.StatusBadRequest)
			// session.AddFlash("Given file is not a valid kubernetes config file, please try again")
			// session.Save(req, w)
			// http.Redirect(w, req, "/play/dashboard", http.StatusFound)
			http.Error(w, "Given file is not a valid kubernetes config file, please try again", http.StatusBadRequest)
			return
		}
		logrus.Debugf("current context: %s, contexts from config file: %v, clusters: %v", ccfg.CurrentContext, ccfg.Contexts, ccfg.Clusters)
		if contextName != "" {
			k8sCtx, ok := ccfg.Contexts[contextName]
			if !ok || k8sCtx == nil {
				logrus.Errorf("error specified context not found")
				// http.Error(w, "context not valid", http.StatusBadRequest)
				// session.AddFlash("Given context name is not valid, please try again with a valid value")
				// session.Save(req, w)
				// http.Redirect(w, req, "/play/dashboard", http.StatusFound)
				http.Error(w, "Given context name is not valid, please try again with a valid value", http.StatusBadRequest)
				return
			}
			// all good, now set the current context to use
			ccfg.CurrentContext = contextName
		}

		// session, err := h.config.SessionStore.Get(req, h.config.SessionName)
		// if err != nil {
		// 	logrus.Errorf("error getting session: %v", err)
		// 	http.Error(w, "unable to get session", http.StatusUnauthorized)
		// 	return
		// }
		session.Values["k8sContext"] = contextName
		session.Values["k8sConfig"] = k8sConfigBytes
	}
	session.Values["k8sInCluster"] = inClusterConfig

	meshLocationURL := req.FormValue("meshLocationURL")
	logrus.Debugf("meshLocationURL: %s", meshLocationURL)
	session.Values["meshLocationURL"] = meshLocationURL

	err = session.Save(req, w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}

	mClient, err := meshes.CreateClient(req.Context(), k8sConfigBytes, contextName, meshLocationURL)
	if err != nil {
		logrus.Errorf("error creating a mesh client: %v", err)
		// http.Error(w, "unable to create an istio client", http.StatusInternalServerError)
		// session.AddFlash("Unable to connect to the mesh using the given kubernetes config file, please try again")
		// session.Save(req, w)
		// http.Redirect(w, req, "/play/dashboard", http.StatusFound)
		http.Error(w, "Unable to connect to the mesh using the given kubernetes config file, please try again", http.StatusInternalServerError)
		return
	}
	defer mClient.Close()
	respOps, err := mClient.MClient.SupportedOperations(req.Context(), &meshes.SupportedOperationsRequest{})

	// meshClient, err := istio.CreateIstioClientWithK8SConfig(ctx, k8sConfigBytes, contextName)
	// if err != nil {
	// 	logrus.Errorf("error creating an istio client: %v", err)
	// 	// http.Error(w, "unable to create an istio client", http.StatusInternalServerError)
	// 	session.AddFlash("Unable to connect to Istio using the given kubernetes config file, please try again")
	// 	session.Save(req, w)
	// 	http.Redirect(w, req, "/play/dashboard", http.StatusFound)
	// 	return
	// }

	// ops, err := meshClient.Operations(ctx)
	if err != nil {
		logrus.Errorf("error getting operations for the mesh: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}

	meshNameOps, err := mClient.MClient.MeshName(req.Context(), &meshes.MeshNameRequest{})
	if err != nil {
		logrus.Errorf("error getting mesh name: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"Ops":  respOps.Ops,
		"Name": meshNameOps.GetName(),
		// "User": user,
	}

	// err = dashTempl.Execute(w, result)
	// if err != nil {
	// 	logrus.Errorf("error rendering the template for the dashboard: %v", err)
	// 	http.Error(w, "unable to render the page", http.StatusInternalServerError)
	// 	return
	// }
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		logrus.Errorf("error marshalling data: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}
}
